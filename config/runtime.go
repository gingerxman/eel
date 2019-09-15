// Copyright 2018 eel Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config

import (
	"context"
	"net/http"

	"github.com/bitly/go-simplejson"
	"github.com/gingerxman/gorm"
)

type runtime struct {
	NewBusinessContext func(ctx context.Context, request *http.Request, userId int, jwtToken string, RawData *simplejson.Json) context.Context

	DB *gorm.DB //database connection
}

func init() {
	Runtime = new(runtime)
}
