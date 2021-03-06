// Copyright (C) 2017 Google Inc.
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

package main

import (
	"context"
	"flag"
	"path/filepath"

	"github.com/google/gapid/core/app"
	"github.com/google/gapid/core/log"
)

type exportReplayVerb struct{ ExportReplayFlags }

func init() {
	verb := &exportReplayVerb{}
	app.AddVerb(&app.Verb{
		Name:      "export_replay",
		ShortHelp: "Export replay vm instruction and assets.",
		Action:    verb,
	})
}

func (verb *exportReplayVerb) Run(ctx context.Context, flags flag.FlagSet) error {
	if flags.NArg() != 1 {
		app.Usage(ctx, "Exactly one gfx trace file expected, got %d", flags.NArg())
		return nil
	}

	capture, err := filepath.Abs(flags.Arg(0))
	if err != nil {
		log.Errf(ctx, err, "Could not find capture file: %v", flags.Arg(0))
	}

	client, err := getGapis(ctx, verb.Gapis, verb.Gapir)
	if err != nil {
		return log.Err(ctx, err, "Failed to connect to the GAPIS server")
	}
	defer client.Close()

	capturePath, err := client.LoadCapture(ctx, capture)
	if err != nil {
		return log.Err(ctx, err, "Failed to load the capture file")
	}

	device, err := getDevice(ctx, client, capturePath, verb.Gapir)
	if err != nil {
		return err
	}

	if err := client.ExportReplay(ctx, capturePath, device, verb.Out); err != nil {
		return log.Err(ctx, err, "Failed to export replay")
	}
	return nil
}
