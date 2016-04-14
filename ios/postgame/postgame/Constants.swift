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
		var result: String
		if let h = NSUserDefaults.standardUserDefaults().stringForKey("host") {
			result = h
		}
		result = host
		if !result.hasPrefix("http:") {
			result = "http://" + result
		}
		return result
	}
}

enum GameMode: Int {
	case Fun = 1, Survival
}