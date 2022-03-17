/*
 * Copyright (C) 2022 The LinQ Authors
 * This file is part of The LinQ library.
 *
 * The LinQ is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The LinQ is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 * You should have received a copy of the GNU Lesser General Public License
 * along with The LinQ.  If not, see <http://www.gnu.org/licenses/>.
 */

package models

type Snapshot struct {
	ID    int64  `gorm:"primaryKey;autoIncrement"`
	Hash  string `gorm:"not null;index:snapshot_hash;"`
	Bytes string
}

type Block struct {
	ID        int64  `gorm:"primaryKey;autoIncrement"`
	BlockHash string `gorm:"not null;uniqueIndex:block_hashcode;index:block_check,priority:1"`
	Height    uint64 `gorm:"not null;uniqueIndex:block_height,sort:desc;index:block_check,priority:2,sort:desc"`
	Bytes     string
}
