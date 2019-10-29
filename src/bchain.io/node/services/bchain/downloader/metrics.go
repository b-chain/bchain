////////////////////////////////////////////////////////////////////////////////
// Copyright (c) 2018 The bchain-go Authors.
//
// The bchain-go is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.
//
// @File: metrics.go
// @Date: 2018/05/08 18:02:08
////////////////////////////////////////////////////////////////////////////////

// Contains the metrics collected by the downloader.

package downloader

import (
	"bchain.io/utils/metrics"
)

var (
	headerInMeter      = metrics.GetOrRegisterMeter("bchain/downloader/headers/in",metrics.DefaultRegistry)
	headerReqTimer     = metrics.NewRegisteredTimer("bchain/downloader/headers/req",metrics.DefaultRegistry)
	headerDropMeter    = metrics.GetOrRegisterMeter("bchain/downloader/headers/drop",metrics.DefaultRegistry)
	headerTimeoutMeter = metrics.GetOrRegisterMeter("bchain/downloader/headers/timeout",metrics.DefaultRegistry)

	bodyInMeter      = metrics.GetOrRegisterMeter("bchain/downloader/bodies/in",metrics.DefaultRegistry)
	bodyReqTimer     = metrics.NewRegisteredTimer("bchain/downloader/bodies/req",metrics.DefaultRegistry)
	bodyDropMeter    = metrics.GetOrRegisterMeter("bchain/downloader/bodies/drop",metrics.DefaultRegistry)
	bodyTimeoutMeter = metrics.GetOrRegisterMeter("bchain/downloader/bodies/timeout",metrics.DefaultRegistry)

	certificateInMeter      = metrics.GetOrRegisterMeter("bchain/downloader/bodies/in",metrics.DefaultRegistry)
	certificateReqTimer     = metrics.NewRegisteredTimer("bchain/downloader/bodies/req",metrics.DefaultRegistry)
	certificateDropMeter    = metrics.GetOrRegisterMeter("bchain/downloader/bodies/drop",metrics.DefaultRegistry)
	certificateTimeoutMeter = metrics.GetOrRegisterMeter("bchain/downloader/bodies/timeout",metrics.DefaultRegistry)

	receiptInMeter      = metrics.GetOrRegisterMeter("bchain/downloader/receipts/in",metrics.DefaultRegistry)
	receiptReqTimer     = metrics.NewRegisteredTimer("bchain/downloader/receipts/req",metrics.DefaultRegistry)
	receiptDropMeter    = metrics.GetOrRegisterMeter("bchain/downloader/receipts/drop",metrics.DefaultRegistry)
	receiptTimeoutMeter = metrics.GetOrRegisterMeter("bchain/downloader/receipts/timeout",metrics.DefaultRegistry)

	stateInMeter   = metrics.GetOrRegisterMeter("bchain/downloader/states/in",metrics.DefaultRegistry)
	stateDropMeter = metrics.GetOrRegisterMeter("bchain/downloader/states/drop",metrics.DefaultRegistry)

)
