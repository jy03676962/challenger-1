//
//  DeviceButton.swift
//  admin
//
//  Created by tassar on 6/26/16.
//  Copyright © 2016 pulupulu. All rights reserved.
//

import UIKit

class DeviceButton: UIButton {

	var controller: PlayerController? {
		didSet {
			if let c = controller {
				isEnabled = true
				if c.matchID == 0 {
					setBackgroundImage(UIImage(named: "PCAvailable"), for: UIControlState())
				} else {
					setBackgroundImage(UIImage(named: "PCGaming"), for: UIControlState())
				}
				let id: String = c.address.id
				var title: String?
				if c.address.type == .simulator {
					title = id[0]
				} else if c.address.type == .wearable {
					title = id.last()
				}
				setTitle(title, for: UIControlState())
				setTitle(title, for: .selected)
			} else {
				isEnabled = false
				isSelected = false
				setTitle(nil, for: UIControlState())
			}
		}
	}
}
