//
//  WebsiteModels.swift
//  postgame
//
//  Created by tassar on 5/10/16.
//  Copyright © 2016 pulupulu. All rights reserved.
//

import Foundation
import ObjectMapper

class BaseResult: Mappable {
	var code: Int?
	var error: String?
	required init?(_ map: Map) {
	}

	func mapping(map: Map) {
		code <- map["code"]
		error <- map["error"]
	}
}

class LoginResult: BaseResult {
	var username: String!
	var userID: Int!
	var currentCoin: Int!
	required init?(_ map: Map) {
		super.init(map)
	}

	override func mapping(map: Map) {
		super.mapping(map)
		username <- map["username"]
		userID <- map["user_id"]
		currentCoin <- map["current_coin"]
	}
}

class AddMatchResult: BaseResult {
	var matchID: Int!
	required init?(_ map: Map) {
		super.init(map)
	}
	override func mapping(map: Map) {
		super.mapping(map)
		matchID <- map["match_id"]
	}
}
