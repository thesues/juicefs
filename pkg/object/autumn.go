/*
 * JuiceFS, Copyright (C) 2018 Juicedata, Inc.
 *
 * This program is free software: you can use, redistribute, and/or modify
 * it under the terms of the GNU Affero General Public License, version 3
 * or later ("AGPL"), as published by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful, but WITHOUT
 * ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or
 * FITNESS FOR A PARTICULAR PURPOSE.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program. If not, see <http://www.gnu.org/licenses/>.
 */

package object

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"time"

	"github.com/journeymidnight/autumn/autumn_clientv1"
)

type autumnStore struct {
	DefaultObjectStorage
	autumnLib *autumn_clientv1.AutumnLib
}

func (d *autumnStore) String() string {
	return "autumn://"
}

func (d *autumnStore) Head(key string) (Object, error) {
	_key, size, err := d.autumnLib.Head(context.Background(), []byte(key))
	if err != nil {
		return nil, err
	}
	return &obj{
		key:   string(_key),
		size:  int64(size),
		isDir: strings.HasSuffix(key, "/"),
		mtime: time.Now(),
	}, nil
}

func (d *autumnStore) List(prefix, marker string, limit int64) ([]Object, error) {
	keys, _, err := d.autumnLib.Range(context.Background(), []byte(prefix), []byte(marker), uint32(limit))
	var ret []Object
	if err != nil {
		return nil, err
	}
	for _, key := range keys {
		ret = append(ret, &obj{
			key:   string(key),
			isDir: strings.HasSuffix(string(key), "/"),
		})
	}
	return ret, nil
}

func (d *autumnStore) Get(key string, off, limit int64) (io.ReadCloser, error) {
	value, err := d.autumnLib.Get(context.Background(), []byte(key))
	if err != nil {
		return nil, err
	}
	return ioutil.NopCloser(bytes.NewReader(value)), nil

}

func (d *autumnStore) Delete(key string) error {
	return d.autumnLib.Delete(context.Background(), []byte(key))
}

func (d *autumnStore) Put(key string, in io.Reader) error {
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(in); err != nil {
		return err
	}
	return d.autumnLib.Put(context.Background(), []byte(key), buf.Bytes())
}

func createAutumnStore(endpoint, accessKey, secretKey string) (ObjectStorage, error) {
	fmt.Printf("connecting..\n")
	client := autumn_clientv1.NewAutumnLib([]string{"127.0.0.1:2379"})
	if err := client.Connect(); err != nil {
		return nil, err
	}
	return &autumnStore{
		autumnLib: client,
	}, nil
}

func init() {
	Register("autumn", createAutumnStore)
}
