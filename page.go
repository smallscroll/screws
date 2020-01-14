package screws

//IPage 分页器接口
type IPage interface {
	Split()
}

//Page 分页器
type Page struct {
	Count           int //总记录数
	PageSize        int //每页显示记录数
	PageCount       int //总页数
	PageIndex       int //当前页码
	LimitStart      int //分页数据截取开始
	LimitEnd        int //分页数据截取结束
	PrePage         int //上一页页码
	NextPage        int //下一页页码
	SelectPageIndex int //请求的页码
}

//NewPage 初始化分页器(总记录数，每页显示记录数)
func NewPage(count int, pageSize int) IPage {
	return &Page{
		Count:           count,
		PageSize:        pageSize,
		PageIndex:       1,
		SelectPageIndex: 0,
	}
}

//Split 分页处理
func (p *Page) Split() {

	if p.Count == 0 || p.Count < p.PageSize {
		p.PageCount = 1
	} else if p.Count > p.PageSize && p.Count%p.PageSize != 0 {
		p.PageCount = p.Count/p.PageSize + 1
	} else {
		p.PageCount = p.Count / p.PageSize
	}

	if p.SelectPageIndex == 0 {
		p.PageIndex = 1
	} else {
		p.PageIndex = p.SelectPageIndex
		if p.PageIndex > p.PageCount {
			p.PageIndex = p.PageCount
		}
		if p.PageIndex < 1 {
			p.PageIndex = 1
		}
	}

	p.LimitStart = (p.PageIndex - 1) * p.PageSize
	if p.Count%p.PageSize == 0 {
		p.LimitEnd = p.LimitStart + p.PageSize
	} else if p.PageIndex < p.PageCount {
		p.LimitEnd = p.LimitStart + p.PageSize
	} else if p.PageIndex >= p.PageCount {
		p.LimitEnd = p.LimitStart + p.Count%p.PageSize
	}
	if p.Count == 0 {
		p.LimitEnd = 0
	}

	if p.PageIndex == 1 {
		p.PrePage = p.PageIndex
	} else {
		p.PrePage = p.PageIndex - 1
	}
	if p.PageIndex == p.PageCount {
		p.NextPage = p.PageCount
	} else {
		p.NextPage = p.PageIndex + 1
	}
}
