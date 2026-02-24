package service

import (
	"context"
	"time"

	"github.com/capyflow/housekeeper/internal/model"
	"github.com/capyflow/housekeeper/internal/repository"
	"github.com/capyflow/housekeeper/pkg"
)

type NotesService struct {
	ctx      context.Context
	noteRepo *repository.NoteRepo
}

func NewNotesService(ctx context.Context, noteRepo *repository.NoteRepo) *NotesService {
	return &NotesService{ctx: ctx, noteRepo: noteRepo}
}

// 创建共享看板
func (ns *NotesService) CreateShareBoard(ctx context.Context, req *model.CreateShareBoardReq) (*model.ShareBoard, error) {
	shareBoardId := "sb_" + pkg.GenerateRandomString(15)
	info := &model.ShareBoard{
		Id:         shareBoardId,
		DeviceName: req.DeviceName,
		BoardName:  req.BoardName,
		Content:    req.Content,
		Owner:      req.Owner,
		CreateTime: time.Now().UnixMilli(),
		UpdateTime: time.Now().UnixMilli(),
	}
	return info, ns.noteRepo.CreateShareBoard(ctx, info)
}

// 更新共享看板
func (ns *NotesService) UpdateShareBoard(ctx context.Context, req *model.UpdateShareBoardReq) error {
	return ns.noteRepo.UpdateShareBoard(ctx, req)
}

// 删除共享看板
func (ns *NotesService) DeleteShareBoard(ctx context.Context, id string) error {
	return ns.noteRepo.DeleteShareBoard(ctx, id)
}

// 分页查询共享看板
func (ns *NotesService) ListShareBoard(ctx context.Context, req *model.ListShareBoardReq) (*model.ListShareBoardResp, error) {
	// 参数校验和默认值设置
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	if req.PageSize > 100 {
		req.PageSize = 100 // 限制最大每页数量
	}

	list, total, err := ns.noteRepo.ListShareBoard(ctx, req)
	if err != nil {
		return nil, err
	}

	return &model.ListShareBoardResp{
		Total: total,
		Page:  req.Page,
		List:  list,
	}, nil
}

// 创建笔记
func (ns *NotesService) CreateNote(ctx context.Context, req *model.CreateNoteReq) (*model.Note, error) {
	noteId := "n_" + pkg.GenerateRandomString(15)
	info := &model.Note{
		Id:       noteId,
		Title:    req.Title,
		Content:  req.Content,
		Owner:    req.Owner,
		Cover:    req.Cover,
		CreateTs: time.Now().UnixMilli(),
		UpdateTs: time.Now().UnixMilli(),
		ModifyTs: time.Now().UnixMilli(),
	}
	return info, ns.noteRepo.CreateNote(ctx, info)
}

// 更新笔记
func (ns *NotesService) UpdateNote(ctx context.Context, req *model.UpdateNoteReq) error {
	return ns.noteRepo.UpdateNote(ctx, req)
}

// 删除笔记
func (ns *NotesService) DeleteNote(ctx context.Context, id string) error {
	return ns.noteRepo.DeleteNote(ctx, id)
}

// 查询笔记详情
func (ns *NotesService) GetNote(ctx context.Context, id string) (*model.Note, error) {
	return ns.noteRepo.GetNote(ctx, id)
}

// 分页查询笔记
func (ns *NotesService) ListNote(ctx context.Context, req *model.ListNoteReq) (*model.ListNoteResp, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	if req.PageSize > 100 {
		req.PageSize = 100
	}

	list, total, err := ns.noteRepo.ListNotes(ctx, req)
	if err != nil {
		return nil, err
	}

	return &model.ListNoteResp{
		Total: total,
		Page:  req.Page,
		List:  list,
	}, nil
}
