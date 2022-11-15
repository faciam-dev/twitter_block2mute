package entity

import (
	"sort"
	"time"
)

type Mute struct {
	id              uint
	userID          uint
	targetTwitterID string
	flag            int
	createdAt       time.Time
	updatedAt       time.Time
}

type Mutes []Mute

func NewMute(
	userID uint,
	targetTwitterID string,
	flag int,
	createdAt time.Time,
	updatedAt time.Time,
) *Mute {
	mute := Mute{}
	mute.userID = userID
	mute.targetTwitterID = targetTwitterID
	mute.flag = flag

	return &mute
}

// getter
func (m *Mute) GetID() uint {
	return m.id
}

func (m *Mute) GetUserID() uint {
	return m.userID
}

func (m *Mute) GetTargetTwitterID() string {
	return m.targetTwitterID
}

func (m *Mute) GetFlag() int {
	return m.flag
}

func (m *Mute) GetCreatedAt() time.Time {
	return m.createdAt
}

func (m *Mute) GetUpdatedAt() time.Time {
	return m.updatedAt
}

// ステータス変更など
// Update
func (m *Mute) Update(
	id uint,
	userID uint,
	targetTwitterID string,
	flag int,
	createdAt time.Time,
	updatedAt time.Time,

) {
	m.id = id
	m.userID = userID
	m.targetTwitterID = targetTwitterID
	m.flag = flag
	m.createdAt = createdAt
	m.updatedAt = updatedAt
}

// Flag
func (m *Mute) Converted() {
	m.flag = 1
}

func (m *Mute) NotConverted() {
	m.flag = 0
}

func (m *Mute) IsConverted() bool {
	return m.flag != 0
}

// ソート
func (ms *Mutes) SortByTargetTwtitterID() {
	sort.Slice(*ms, func(i, j int) bool {
		return (*ms)[i].GetTargetTwitterID() <= (*ms)[j].GetTargetTwitterID()
	})
}

// 検索
// 対象のMutesに存在しないIDsを検索する。
func (ms *Mutes) IsConvertedByTwitterID(needleTwitterID string) bool {
	idx := sort.Search(len(*ms), func(i int) bool {
		return string((*ms)[i].GetTargetTwitterID()) == needleTwitterID
	})

	if len(*ms) > idx && (*ms)[idx].IsConverted() {
		return true
	}

	return false
}

// コンバート
// Mutes への変換
func (ms *Mutes) ToMutesEntity(mutes *[]Mute) {
	for _, v := range *mutes {
		*ms = append(*ms, v)
	}
}

// *[]Mute への変換
func (ms *Mutes) ToMuteEntities() *[]Mute {
	mutes := []Mute{}
	for _, v := range *ms {
		mutes = append(mutes, v)
	}
	return &mutes
}
