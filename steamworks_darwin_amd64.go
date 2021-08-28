// SPDX-License-Identifier: Apache-2.0
// Copyright 2021 The go-steamworks Authors

package steamworks

import (
	_ "embed"
)

//go:embed libsteam_api.dylib
var libSteamAPI []byte
