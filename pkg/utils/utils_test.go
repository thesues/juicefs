/*
 * JuiceFS, Copyright (C) 2021 Juicedata, Inc.
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

package utils

import (
	"testing"
	"time"
)

func TestMin(t *testing.T) {
	assertEqual(t, Min(1, 2), 1)
	assertEqual(t, Min(-1, -2), -2)
	assertEqual(t, Min(0, 0), 0)
}

func TestExists(t *testing.T) {
	assertEqual(t, Exists("/"), true)
	assertEqual(t, Exists("/not_exist_path"), false)
}

func TestSplitDir(t *testing.T) {
	assertEqual(t, SplitDir("/a:/b"), []string{"/a", "/b"})
	assertEqual(t, SplitDir("a,/b"), []string{"a", "/b"})
	assertEqual(t, SplitDir("/a;b"), []string{"/a;b"})
	assertEqual(t, SplitDir("a/b"), []string{"a/b"})
}

func TestGetInode(t *testing.T) {
	_, err := GetFileInode("")
	if err == nil {
		t.Fatalf("invalid path should fail")
	}
	ino, err := GetFileInode("/")
	if err != nil {
		t.Fatalf("get file inode: %s", err)
	} else if ino > 2 {
		t.Fatalf("inode of root should be 1/2, but got %d", ino)
	}
}

func TestProgresBar(t *testing.T) {
	p, bar := NewProgressCounter("test")
	go func() {
		for i := 0; i < 100; i++ {
			time.Sleep(time.Millisecond)
			bar.Increment()
		}
		bar.SetTotal(0, true)
	}()
	p.Wait()

	p, bar = NewDynProgressBar("test", true)
	go func() {
		for i := 0; i < 100; i++ {
			time.Sleep(time.Millisecond)
			bar.Increment()
		}
		bar.SetTotal(0, true)
	}()
	p.Wait()
}

func TestLocalIp(t *testing.T) {
	_, err := GetLocalIp("127.0.0.1")
	if err == nil {
		t.Fatalf("should fail with invalid address")
	}
	ip, err := GetLocalIp("127.0.0.1:22")
	if err != nil {
		t.Fatalf("get local ip: %s", err)
	}
	if ip != "127.0.0.1" {
		t.Fatalf("local ip should be 127.0.0.1, bug got %s", ip)
	}
}
