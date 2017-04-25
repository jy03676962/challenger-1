//
//  ArduinoController.swift
//  admin
//
//  Created by tassar on 5/5/16.
//  Copyright © 2016 pulupulu. All rights reserved.
//

import Foundation
import ObjectMapper

enum ArduinoMode: Int {
	case unknown = -1, off = 0, on, scan, free
}

class ArduinoController: Mappable {

	var address: Address!
	var id: String!
	var mode: ArduinoMode!
	var online: Bool!
	var scoreUpdated: Bool!

	required init?(map: Map) {
	}

	func mapping(map: Map) {
		address <- map["address"]
		id <- map["id"]
		mode <- map["mode"]
		online <- map["online"]
		scoreUpdated <- map["scoreUpdated"]
	}
}
