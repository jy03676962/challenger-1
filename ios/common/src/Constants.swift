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

/**
 @brief 不要用这里的做为ViewTag!!
 */
struct ReservedViewTag {
	static let vc_BackgroundImageView = 10001
}

enum GameMode: Int {
	case Fun = 1, Survival
}