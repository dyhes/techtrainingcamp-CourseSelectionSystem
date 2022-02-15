package course_cache

import (
	"github.com/dyhes/techtrainingcamp-CourseSelectionSystem/models"
	"github.com/dyhes/techtrainingcamp-CourseSelectionSystem/pkg/e"
	"sync"
)

/*
课程信息缓存
*/
type CourseInfoCache struct {
	info map[uint64]*models.CourseInfo //key为courseID，value为CourseInfo
	lock sync.RWMutex
}

//添加课程信息缓存
func (c *CourseInfoCache) Add(courseID uint64) {
	c.lock.Lock()
	defer c.lock.Unlock()
	info, _ := models.GetCourseInfo(courseID)
	c.info[courseID] = &info
}

//获取某一课程的信息缓存
func (c *CourseInfoCache) Get(courseID uint64) *models.CourseInfo {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.info[courseID]
}

//批量获取课程信息缓存
func (c *CourseInfoCache) GetList(ids []uint64) []models.CourseInfo {
	c.lock.RLock()
	defer c.lock.RUnlock()
	n := len(ids)
	res := make([]models.CourseInfo, n)
	for i := 0; i < n; i++ {
		res[i] = *c.info[ids[i]]
	}
	return res
}

/*
课程缓存
*/
type courseCache struct {
	//key为courseID，value为studentIDSet，即选了这门课的所有学生的ID
	course map[uint64]*IDSet

	//key为courseID，value为该课程的剩余容量
	remainCap map[uint64]*Cap

	//用于添加课程的读写锁
	cLock sync.RWMutex

	//用于修改课程剩余容量的互斥锁
	rLock sync.RWMutex
}

//查找课程缓存中是否存在某课程
func (c *courseCache) IsCourseExist(courseID uint64) bool {
	c.cLock.RLock()
	defer c.cLock.RUnlock()
	_, ok := c.course[courseID]
	return ok
}

//向课程缓存中添加课程。不是线程安全且默认课程缓存中不存在将要添加的课程
func (c *courseCache) AddCourse(courseID uint64, cap uint) {
	c.course[courseID] = &IDSet{set: map[uint64]bool{}}
	c.remainCap[courseID] = &Cap{remainCap: cap}

	InfoCache.Add(courseID)
}

//添加选择了某一课程的学生ID。courseID和studentID必须都在缓存中
func (c *courseCache) AddStudent(courseID, studentID uint64) e.ErrNo {
	//如果学生已经选择了该课程，返回e.StudentHasCourse
	if StudentCache.HasCourse(studentID, courseID) {
		return e.StudentHasCourse
	}

	//若课程没有名额了，返回e.CourseNotAvailable
	c.rLock.RLock()
	if c.remainCap[courseID].Zero() {
		c.rLock.RUnlock()
		return e.CourseNotAvailable
	}
	c.rLock.RUnlock()

	//只有在学生没有选择该课程且该课程还有名额时，才添加课程
	if !StudentCache.AddCourse(studentID, courseID) {
		return e.StudentHasCourse
	}

	//课程容量减1
	c.rLock.RLock()
	if !c.remainCap[courseID].MinusOne() {
		//课程容量已满，所有要删除刚才添加的课程
		c.rLock.RUnlock()
		StudentCache.DeleteCourse(studentID, courseID)
		return e.CourseNotAvailable
	}
	c.rLock.RUnlock()

	c.cLock.RLock()
	c.course[courseID].AddID(studentID)
	c.cLock.RUnlock()

	return e.OK
}

/*
学生缓存
*/
type studentCache struct {
	//key为studentID，value为courseIDSet，即该学生选择的所有课程的ID
	student map[uint64]*IDSet

	//读写锁
	lock sync.RWMutex
}

//查找学生缓存中是否存在某学生
func (s *studentCache) IsStudentExist(studentID uint64) bool {
	s.lock.RLock()
	defer s.lock.RUnlock()
	_, ok := s.student[studentID]
	return ok
}

//向学生缓存中添加学生。不是线程安全且默认学生缓存中不存在将要添加的学生
func (s *studentCache) AddStudent(studentID uint64) {
	s.student[studentID] = &IDSet{set: map[uint64]bool{}}
}

//判断学生是否已经选择了该课程
func (s *studentCache) HasCourse(studentID, courseID uint64) bool {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.student[studentID].ExistID(courseID)
}

//为学生添加课程。若学生已经选择了该课程，返回false。否则添加课程并返回true
func (s *studentCache) AddCourse(studentID, courseID uint64) bool {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.student[studentID].AddID(courseID)
}

//为学生删除课程
func (s *studentCache) DeleteCourse(studentID, courseID uint64) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	s.student[studentID].DeleteID(courseID)
}

/*
ID set集合(学生ID或者课程ID)
*/
type IDSet struct {
	set  map[uint64]bool
	lock sync.RWMutex
}

//添加ID。若ID已存在，返回false，否则添加并返回true
func (i *IDSet) AddID(id uint64) bool {
	i.lock.Lock()
	defer i.lock.Unlock()
	if _, ok := i.set[id]; ok {
		//id已存在，返回false
		return false
	}
	i.set[id] = true
	return true
}

//删除ID
func (i *IDSet) DeleteID(id uint64) {
	i.lock.Lock()
	defer i.lock.RUnlock()
	delete(i.set, id)
}

//判断ID是否存在
func (i *IDSet) ExistID(id uint64) bool {
	i.lock.RLock()
	defer i.lock.RUnlock()
	_, ok := i.set[id]
	return ok
}

//返回集合中的所有ID
func (i *IDSet) GetAll() []uint64 {
	i.lock.RLock()
	defer i.lock.RUnlock()
	res := make([]uint64, 0)
	for id := range i.set {
		res = append(res, id)
	}
	return res
}

/*
课程容量
*/
type Cap struct {
	remainCap uint
	lock      sync.RWMutex
}

//课程剩余容量为0时，返回true。否则，返回false
func (c *Cap) Zero() bool {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.remainCap == 0
}

//若课程无剩余名额，返回false。若课程还有剩余名额，剩余名额减一，返回true
func (c *Cap) MinusOne() bool {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.remainCap == 0 {
		return false
	}
	c.remainCap -= 1
	return true
}

/*
同步
*/
type OnceFlag struct {
	flag map[uint64]*sync.Once
	sig  map[uint64]chan struct{}
	lock sync.RWMutex
}

func (o *OnceFlag) Exist(studentID uint64) bool {
	o.lock.RLock()
	defer o.lock.RUnlock()
	_, ok := o.flag[studentID]
	return ok
}

func (o *OnceFlag) Add(studentID uint64) {
	o.lock.Lock()
	defer o.lock.Unlock()
	_, ok := o.flag[studentID]
	if ok {
		return
	}
	o.flag[studentID] = &sync.Once{}
	o.sig[studentID] = make(chan struct{})
}

func (o *OnceFlag) Get(studentID uint64) (*sync.Once, chan struct{}) {
	o.lock.RLock()
	defer o.lock.RUnlock()
	return o.flag[studentID], o.sig[studentID]
}

var (
	//学生缓存
	StudentCache = studentCache{student: map[uint64]*IDSet{}}
	//课程缓存
	CourseCache = courseCache{course: map[uint64]*IDSet{}, remainCap: map[uint64]*Cap{}}
	//课程信息缓存key为courseID，value为CourseInfo
	InfoCache = CourseInfoCache{info: map[uint64]*models.CourseInfo{}}

	//保证将学生添加到学生缓存中的操作只进行一次
	OnceStudent = OnceFlag{flag: map[uint64]*sync.Once{}, sig: map[uint64]chan struct{}{}}
	//保证将课程添加到课程缓存中的操作只进行一次
	OnceCourse = OnceFlag{flag: map[uint64]*sync.Once{}, sig: map[uint64]chan struct{}{}}
)

//判断studentID是否存在。对同一个studentID只会查询一次数据库。
func CheckStudentExistence(studentID uint64) bool {
	OnceStudent.Add(studentID)
	o, ch := OnceStudent.Get(studentID)
	o.Do(func() {
		defer close(ch)
		//首先到缓存中查找
		if !StudentCache.IsStudentExist(studentID) {
			//缓存中不存在时去数据库中查找
			exist, err := models.IsUserExistByID(studentID, 2)
			if err != nil {
				return
			}
			if !exist {
				return
			}
			//数据库中存在，则添加到缓存
			StudentCache.AddStudent(studentID)
		}
	})

	<-ch

	//再查询一次缓存。如果这里还没有查询到，说明studentID不存在
	if !StudentCache.IsStudentExist(studentID) {
		return false
	}
	return true
}

//判断courseID是否存在。对同一个courseID只会查询一次数据库
func CheckCourseExistence(courseID uint64) bool {
	OnceCourse.Add(courseID)
	o, ch := OnceCourse.Get(courseID)
	o.Do(func() {
		defer close(ch)
		//首先查询缓存
		if !CourseCache.IsCourseExist(courseID) {
			//缓存不存在再查询数据库
			exist, _ := models.IsCourseExistByID(courseID)
			if !exist {
				return
			}
			//添加课程到缓存
			rCap, _ := models.GetCourseRemainCap(courseID)
			CourseCache.AddCourse(courseID, rCap)
		}

	})

	//所有courseID相同的请求都必须要等待channel关闭。
	<-ch

	//channel关闭时，再查询一次缓存。如果这里还没有查询到，说明courseID不存在
	if !CourseCache.IsCourseExist(courseID) {
		return false
	}
	return true
}

//学生选择课程
func StudentSelectCourse(studentID, courseID uint64) e.ErrNo {
	return CourseCache.AddStudent(courseID, studentID)
}

//获取学生课程表
func GetStudentCourse(studentID uint64) []models.CourseInfo {
	StudentCache.lock.RLock()
	coursesID := StudentCache.student[studentID].GetAll()
	StudentCache.lock.RUnlock()

	return InfoCache.GetList(coursesID)
}

//将缓存中的信息写入数据库
func WriteToDatabase() {

}
