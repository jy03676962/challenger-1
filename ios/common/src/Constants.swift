//
//  Constants.swift
//  postgame
//
//  Created by tassar on 4/7/16.
//  Copyright © 2016 pulupulu. All rights reserved.
//

import Foundation
import SwiftyUserDefaults

struct PLConstants {
	static let host = "localhost:3030"
	static let usualFont = "Alien League Bold"
	static let maxTeamSize = 4
	static func getHost() -> String {
		if let h = Defaults[.host] {
			return h
		}
		return host
	}
	static func getClientWsAddress() -> String {
		return "ws://" + getHost() + "/client"
	}
	static func getAdminWsAddress() -> String {
		return "ws://" + getHost() + "/admin"
	}
}

enum GameMode: Int {
	case Fun = 1, Survival
}