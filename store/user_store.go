package store

import (
	"context"
	"net/http"
	"github.com/SekiguchiKai/learn_GAE_Go_PropertyLoadSaver/model"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

const _UserKind = "User"

type UserStore struct {
	ctx context.Context
}

func NewUserStore(r *http.Request) UserStore {
	return NewUserStoreWithContext(appengine.NewContext(r))
}

func NewUserStoreWithContext(ctx context.Context) UserStore {
	return UserStore{ctx: ctx}
}

//
func (s UserStore) GetUserList(cursor string, limit int, dst *model.UserList) error {
	// Query発行
	q := datastore.NewQuery(_UserKind).Limit(limit)

	// CursorをDecode
	if start, err := datastore.DecodeCursor(cursor); err == nil {
		// cursorを元に途中の位置から始まるQueryを発行する
		q = q.Start(start)
	}

	// 格納用のslice
	var users []model.User

	// contextを元にQueryをrunする
	iterator := q.Run(s.ctx)
	for {
		var u model.User
		// これ以上存在しない場合は、datastore.Doneを吐く
		// そうでない場合は、引数のdstにEntityをloadする
		if _, err := iterator.Next(&u); err != nil {
			if err == datastore.Done {
				break
			} else {
				return err
			}
		}
		// sliceに1個分格納する
		users = append(users, u)
	}

	// イテレータの現在の場所のカーソルを返す
	nextCursor, err := iterator.Cursor()
	if err != nil {
		return err
	}

	dst.List = users
	if s.hasNext(q, nextCursor) {
		dst.HasNext = true
		dst.Cursor = nextCursor.String()
	}

	return nil
}

func (s UserStore) GetUser(id string, dst *model.User) (exists bool, e error) {
	if id == "" {
		return false, nil
	}

	key := s.newUserKey(id)
	if err := datastore.Get(s.ctx, key, dst); err != nil {
		if err != datastore.ErrNoSuchEntity {
			return false, err
		}
		return false, nil
	}
	return true, nil
}


func (s UserStore) ExistsUser(id string) (bool, error) {
	var dst model.User
	return s.GetUser(id, &dst)
}

func (s UserStore) PutUser(u model.User) error {
	key := s.newUserKey(u.ID)
	if _, err := datastore.Put(s.ctx, key, &u); err != nil {
		return err
	}
	return nil
}


func (s UserStore) DeleteUser(id string) error {
	key := s.newUserKey(id)
	return datastore.Delete(s.ctx, key)
}

func (s UserStore) hasNext(q *datastore.Query, c datastore.Cursor) bool {
	if _, err := q.Limit(1).Start(c).Run(s.ctx).Next(nil); err == nil {
		return true
	}
	return false
}

func (s UserStore) newUserKey(id string) *datastore.Key {
	return datastore.NewKey(s.ctx, _UserKind, id, 0, nil)
}