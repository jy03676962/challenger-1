//
//  Constants.swift
//  postgame
//
//  Created by tassar on 4/7/16.
//  Copyright © 2016 pulupulu. All rights reserved.
//

import Foundation

struct PLConstants {
	static let host = "localhost:3030"
	static let usualFont = "Alien League Bold"
	static func getHost() -> String {
		if let h = NSUserDefaults.standardUserDefaults().stringForKey("host") {
			return h
		}
		return host
	}
	static func getWsAddress() -> String {
		return "ws://" + getHost() + "/client"
	}
}

enum GameMode: Int {
	case Fun = 1, Survival
}