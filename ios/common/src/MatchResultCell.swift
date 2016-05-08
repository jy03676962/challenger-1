//
//  MatchResultCell.swift
//  admin
//
//  Created by tassar on 5/8/16.
//  Copyright © 2016 pulupulu. All rights reserved.
//

import UIKit

class MatchResultCell: UITableViewCell {
	@IBOutlet weak var idLabel: UILabel!
	@IBOutlet weak var levelLabel: UILabel!
	@IBOutlet weak var gradeLabel: UILabel!
	@IBOutlet weak var goldLabel: UILabel!
	@IBOutlet weak var energyLabel: UILabel!
	@IBOutlet weak var comboLabel: UILabel!

	override func awakeFromNib() {
		super.awakeFromNib()
		backgroundColor = UIColor.clearColor()
	}

	func setData(data: PlayerData) {
		idLabel.text = data.getName()
		levelLabel.text = "LEVEL.\(data.level)"
		gradeLabel.text = data.grade.uppercaseString
		goldLabel.text = "\(data.gold)/-\(data.lostGold)"
		energyLabel.text = "\(Int(data.energy))"
		comboLabel.text = "\(data.combo)"
	}
}
