package utility

var (
	maxNum  int
	nx, ny  int     //x集合和y集合中顶点的数量
	edge    [][]int //邻接链表。edge[i][j]表示x[i]和y[j]是一条边
	visited []bool  //已访问的节点visited[i]为true表示i已访问，为false则表示未访问
	cx, cy  []int   //最大匹配的结果。x集合的第i个元素匹配的是y集合的第cx[i]个元素
)

//生成邻接链表
func MaxMatch(data map[string][]string) map[string]string {
	//先用一个set把所有的课程ID保存下来
	courseVec := make([]string, 0)     //课程ID集合
	courseSet := make(map[string]int)  //课程ID以及它在courseVec中的索引
	teacherVec := make([]string, 0)    //教师ID集合
	teacherSet := make(map[string]int) //教师ID以及它在teacherVec中的索引
	ci, ti := 0, 0
	for teacher, courses := range data {
		//添加教师ID
		teacherVec = append(teacherVec, teacher)
		teacherSet[teacher] = ti
		ti++
		for _, c := range courses {
			if _, ok := courseSet[c]; !ok {
				//添加课程ID
				courseVec = append(courseVec, c)
				courseSet[c] = ci
				ci++
			}
		}
	}

	//生成邻接矩阵
	nx, ny = len(teacherVec), len(courseVec)
	maxNum = max(nx, ny)
	edge = make([][]int, maxNum)
	for i := 0; i < maxNum; i++ {
		edge[i] = make([]int, maxNum)
	}
	for teacher, courses := range data {
		i := teacherSet[teacher]
		for _, course := range courses {
			j := courseSet[course]
			edge[i][j] = 1
		}
	}

	//最大匹配求解
	cx, cy = make([]int, maxNum), make([]int, maxNum)
	for i := 0; i < maxNum; i++ {
		cx[i] = -1
		cy[i] = -1
	}
	for i := 0; i < nx; i++ {
		if cx[i] == -1 {
			visited = make([]bool, maxNum)
			dfs(i)
		}
	}

	//处理结果
	res := make(map[string]string) //key为教师ID，value最终绑定的课程ID
	for i, _ := range cx {
		res[teacherVec[i]] = courseVec[cx[i]]
	}
	clear()
	return res
}

//寻找从u出发的增广路径
func dfs(u int) bool {
	for v := 0; v < ny; v++ {
		if edge[u][v] > 0 && !visited[v] {
			visited[v] = true
			if cy[v] == -1 || dfs(cy[v]) {
				cx[u] = v
				cy[v] = u
				return true
			}
		}
	}
	return false
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func clear() {
	maxNum = 0
	nx, ny = 0, 0
	edge = make([][]int, 0)
	visited = make([]bool, 0)
	cx, cy = make([]int, 0), make([]int, 0)
}
