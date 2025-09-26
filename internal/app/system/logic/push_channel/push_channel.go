package logic

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	v1 "github.com/shichen437/gowlive/api/v1/system"
	"github.com/shichen437/gowlive/internal/app/system/dao"
	"github.com/shichen437/gowlive/internal/app/system/model/do"
	"github.com/shichen437/gowlive/internal/app/system/model/entity"
	"github.com/shichen437/gowlive/internal/app/system/service"
	"github.com/shichen437/gowlive/internal/pkg/utils"
)

type sPushChannel struct {
}

func init() {
	service.RegisterPushChannel(New())
}

func New() service.IPushChannel {
	return &sPushChannel{}
}

func (s *sPushChannel) List(ctx context.Context, req *v1.GetPushChannelListReq) (res *v1.GetPushChannelListRes, err error) {
	res = &v1.GetPushChannelListRes{}
	m := dao.PushChannel.Ctx(ctx)
	if req.Name != "" {
		m = m.WhereLike(dao.PushChannel.Columns().Name, "%"+req.Name+"%")
	}
	if req.Type != "" {
		m = m.Where(dao.PushChannel.Columns().Type, req.Type)
	}
	res.Total, err = m.Count()
	if err != nil || res.Total == 0 {
		return
	}
	m = m.OrderDesc(dao.PushChannel.Columns().Id)
	err = m.Page(req.PageNum, req.PageSize).Scan(&res.Rows)
	if err != nil {
		return nil, gerror.New("获取渠道列表失败")
	}
	return
}

func (s *sPushChannel) Post(ctx context.Context, req *v1.PostPushChannelReq) (res *v1.PostPushChannelRes, err error) {
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		lastInfo, err := dao.PushChannel.Ctx(ctx).Insert(do.PushChannel{
			Name:      req.Name,
			Type:      req.Type,
			Status:    req.Status,
			Url:       req.Url,
			Remark:    req.Remark,
			CreatedAt: utils.Now(),
		})
		if err != nil {
			return err
		}
		channelId, err := lastInfo.LastInsertId()
		if err != nil {
			return err
		}
		if req.Type == "email" {
			_, err = dao.PushChannelEmail.Ctx(ctx).Insert(do.PushChannelEmail{
				ChannelId: channelId,
				Sender:    req.Email.Sender,
				Receiver:  req.Email.Receiver,
				Server:    req.Email.Server,
				Port:      req.Email.Port,
				AuthCode:  req.Email.AuthCode,
				CreatedAt: utils.Now(),
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	return
}

func (s *sPushChannel) Put(ctx context.Context, req *v1.PutPushChannelReq) (res *v1.PutPushChannelRes, err error) {
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		_, err := dao.PushChannel.Ctx(ctx).WherePri(req.Id).Update(do.PushChannel{
			Name:      req.Name,
			Type:      req.Type,
			Status:    req.Status,
			Url:       req.Url,
			Remark:    req.Remark,
			UpdatedAt: utils.Now(),
		})
		if err != nil {
			return err
		}
		if req.Type == "email" {
			_, err = dao.PushChannelEmail.Ctx(ctx).Where(dao.PushChannelEmail.Columns().ChannelId, req.Id).Update(do.PushChannelEmail{
				Sender:    req.Email.Sender,
				Receiver:  req.Email.Receiver,
				Server:    req.Email.Server,
				Port:      req.Email.Port,
				AuthCode:  req.Email.AuthCode,
				UpdatedAt: utils.Now(),
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	return
}

func (s *sPushChannel) Get(ctx context.Context, req *v1.GetPushChannelReq) (res *v1.GetPushChannelRes, err error) {
	res = &v1.GetPushChannelRes{}
	err = dao.PushChannel.Ctx(ctx).WherePri(req.Id).Scan(&res)
	if err != nil || res == nil || res.Id == 0 {
		err = gerror.New("获取渠道详情失败")
		return
	}
	if res.Type == "email" {
		var email *entity.PushChannelEmail
		err = dao.PushChannelEmail.Ctx(ctx).Where(dao.PushChannelEmail.Columns().ChannelId, req.Id).Limit(1).Scan(&email)
		if err != nil {
			err = gerror.New("获取渠道详情失败")
			return
		}
		if email != nil {
			res.Email = email
		}
	}
	return
}

func (s *sPushChannel) Delete(ctx context.Context, req *v1.DeletePushChannelReq) (res *v1.DeletePushChannelRes, err error) {
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		_, err = dao.PushChannel.Ctx(ctx).WherePri(req.Id).Delete()
		if err != nil {
			return err
		}
		_, err = dao.PushChannelEmail.Ctx(ctx).Where(dao.PushChannelEmail.Columns().ChannelId, req.Id).Delete()
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, gerror.New("删除推送渠道失败")
	}
	return
}
