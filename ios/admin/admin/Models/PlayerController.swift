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
	case unknown = 0, admin, simulator, arduinoTest, postgame, wearable, mainArduino, subArduino, queueDevice, ingameDevice, musicArduino, doorArduino
}

struct Address: Mappable {
	var type: AddressType!
	var id: String!
	init?(map: Map) {
	}
	mutating func mapping(map: Map) {
		type <- map["type"]
		id <- map["id"]
	}
}

class PlayerController: Mappable {
	var address: Address!
	var id: String!
	var matchID: Int!
	var online: Bool!

	required init?(map: Map) {
	}

	func mapping(map: Map) {
		address <- map["address"]
		id <- map["id"]
		matchID <- map["matchID"]
		online <- map["online"]
	}
}
