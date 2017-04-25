//
//  Team.swift
//  admin
//
//  Created by tassar on 4/26/16.
//  Copyright © 2016 pulupulu. All rights reserved.
//

import Foundation
import ObjectMapper

enum TeamStatus: Int {
	case waiting = 0, prepare, playing, after, finished
}

class Team: Mappable {
	var size: Int!
	var id: String!
	var delayCount: Int!
	var status: TeamStatus!
	var waitTime: Int!
	var mode: String!

	required init?(map: Map) {
	}

	func mapping(map: Map) {
		size <- map["size"]
		id <- map["id"]
		delayCount <- map["delayCount"]
		status <- map["status"]
		waitTime <- map["waitTime"]
		mode <- map["mode"]
	}
}
