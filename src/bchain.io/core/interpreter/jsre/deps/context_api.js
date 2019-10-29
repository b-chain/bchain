/*
 * Copyright (c) 2018 The bchain-go Authors.
 *
 * The bchain-go is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 * @File: bchainotto.js
 * @Date: 2018/08/06 10:38:06
 */

/*!(function contextApi() {
    this.auth = _context_api_auth
    this.sys = _context_api_sys
    //this.console = _context_api_console
    this.assert = _context_api_assert
    this.crypto = _context_api_crypto
    this.producer = _context_api_producer
    this.act = _context_api_act
    this.db = _context_api_db
})();*/

;(function () {
    this.ctxApi = function() {
        return {
            _contract: _context_var_self,
            _sender: _context_var_sender,
            _creator:_context_var_creator,
            _miner : _context_var_miner,
            _number: _context_var_number,

            auth: _context_api_auth,
            sys: _context_api_sys,
            console: _context_api_console,
            assert: _context_api_assert,
            crypto: _context_api_crypto,
            producer: _context_api_producer,
            act: _context_api_act,
            db: _context_api_db,
            cache:_context_api_cache,
            contract: _context_api_contract,
            call: _context_api_call,
            result:_context_api_result
        };
    };
})(this);
