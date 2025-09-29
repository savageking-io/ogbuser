package group

import "sync"

type GroupsData struct {
	groups map[int32]*Group
	mutex  sync.Mutex
}

func NewGroupsData() *GroupsData {
	return &GroupsData{
		groups: make(map[int32]*Group),
	}
}

func (d *GroupsData) Add(group *Group) {
	defer d.mutex.Unlock()
	d.mutex.Lock()
	d.groups[int32(group.raw.Id)] = group
}

func (d *GroupsData) Remove(id int32) {
	defer d.mutex.Unlock()
	d.mutex.Lock()
	delete(d.groups, id)
}

func (d *GroupsData) Get(id int32) (*Group, bool) {
	defer d.mutex.Unlock()
	d.mutex.Lock()
	group, ok := d.groups[id]
	return group, ok
}

func (d *GroupsData) GetAll() []*Group {
	var result []*Group
	for _, group := range d.groups {
		result = append(result, group)
	}
	return result
}
