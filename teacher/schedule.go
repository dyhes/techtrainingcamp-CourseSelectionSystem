package teacher

import (
	"net/http"

	"github.com/dyhes/techtrainingcamp-CourseSelectionSystem/pkg/types"
	"github.com/gin-gonic/gin"
)

type match struct {
	edge      [][]int
	n, m      int
	pair, vis []int
}

func (t *match) addEdge(x, y int) {
	t.edge[x] = append(t.edge[x], y)
}

func (t *match) could(rt int) bool {
	for _, to := range t.edge[rt] {
		if t.vis[to] == 0 {
			t.vis[to] = 1
			if t.pair[to] == -1 || t.could(t.pair[to]) {
				t.pair[to] = rt
				return true
			}
		}
	}
	return false
}

func (t *match) calc() int {
	var ans int = 0
	for i := 0; i < t.n; i++ {
		t.vis = make([]int, t.m)
		if t.could(i) {
			ans++
		}
	}
	return ans
}

func (t *match) befor() {
	t.edge = make([][]int, t.n)
	t.pair = make([]int, t.m)
	for i := 0; i < t.m; i++ {
		t.pair[i] = -1
	}
}

func Schedule(c *gin.Context) { // /api/v1/course/schedule
	var req types.ScheduleCourseRequest
	c.BindJSON(&req)

	var n, m int = 0, 0
	teacher := make(map[string]int)
	course := make(map[string]int)
	var teachers, courses []string
	for key, val := range req.TeacherCourseRelationShip {
		if _, have := teacher[key]; !have { // 如果没有出现的老师，就新创建一个id
			teacher[key] = n
			teachers = append(teachers, key)
			n++
		}
		for _, it := range val {
			if _, have := course[it]; !have { // 如果没有出现的课程，就新创建一个id
				course[it] = m
				courses = append(courses, it)
				m++
			}
		}
	}

	var calc match
	calc.n, calc.m = n, m
	calc.befor()

	for u, val := range req.TeacherCourseRelationShip {
		for _, v := range val {
			calc.addEdge(teacher[u], course[v])
		}
	}

	calc.calc()

	var res types.ScheduleCourseResponse
	res.Data = make(map[string]string)

	for i := 0; i < calc.m; i++ {
		if calc.pair[i] != -1 {
			res.Data[teachers[calc.pair[i]]] = courses[i]
		}
	}

	c.JSON(http.StatusOK, res)
}

// type ScheduleCourseRequest struct {
// 	TeacherCourseRelationShip map[string][]string // key 为 teacherID , val 为老师期望绑定的课程 courseID 数组
// }

// type ScheduleCourseResponse struct {
// 	Code ErrNo
// 	Data map[string]string // key 为 teacherID , val 为老师最终绑定的课程 courseID
// }
