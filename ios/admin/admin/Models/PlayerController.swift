//
//  Address.swift
//  admin
//
//  Created by tassar on 5/1/16.
//  Copyright © 2016 pulupulu. All rights reserved.
//

import Foundation
import ObjectMapper

enum AddressType: Int {
	case Unknown = 0, Admin, Simulator, ArduinoTest, Postgame, Wearable, Arduino
}

struct Address: Mappable {
	var type: AddressType!
	var id: String!
	init?(_ map: Map) {
	}
	mutating func mapping(map: Map) {
		type <- map["type"]
		id <- map["id"]
	}
}

enum PCStatus: Int {
	case Offline = 0, Idle, Using
}

struct PlayerController: Mappable {
	var address: Address!
	var status: PCStatus!
	var id: String!

	init?(_ map: Map) {
	}

	mutating func mapping(map: Map) {
		address <- map["address"]
		status <- map["status"]
		id <- map["id"]
	}
}