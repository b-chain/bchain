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

package fetcher

import "bchain.io/utils/metrics"

var (
	propAnnounceInMeter   = metrics.NewRegisteredMeter("bchain/fetcher/prop/announces/in",nil)
	propAnnounceOutTimer  = metrics.NewRegisteredTimer("bchain/fetcher/prop/announces/out",nil)
	propAnnounceDropMeter = metrics.NewRegisteredMeter("bchain/fetcher/prop/announces/drop",nil)
	propAnnounceDOSMeter  = metrics.NewRegisteredMeter("bchain/fetcher/prop/announces/dos",nil)

	propBroadcastInMeter   = metrics.NewRegisteredMeter("bchain/fetcher/prop/broadcasts/in",nil)
	propBroadcastOutTimer  = metrics.NewRegisteredTimer("bchain/fetcher/prop/broadcasts/out",nil)
	propBroadcastDropMeter = metrics.NewRegisteredMeter("bchain/fetcher/prop/broadcasts/drop",nil)
	propBroadcastDOSMeter  = metrics.NewRegisteredMeter("bchain/fetcher/prop/broadcasts/dos",nil)

	headerFetchMeter = metrics.NewRegisteredMeter("bchain/fetcher/fetch/headers",nil)
	bodyFetchMeter   = metrics.NewRegisteredMeter("bchain/fetcher/fetch/bodies",nil)

	headerFilterInMeter  = metrics.NewRegisteredMeter("bchain/fetcher/filter/headers/in",nil)
	headerFilterOutMeter = metrics.NewRegisteredMeter("bchain/fetcher/filter/headers/out",nil)
	bodyFilterInMeter    = metrics.NewRegisteredMeter("bchain/fetcher/filter/bodies/in",nil)
	bodyFilterOutMeter   = metrics.NewRegisteredMeter("bchain/fetcher/filter/bodies/out",nil)
)