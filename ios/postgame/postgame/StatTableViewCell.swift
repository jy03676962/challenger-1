//
//  StatTableViewCell.swift
//  postgame
//
//  Created by tassar on 4/13/16.
//  Copyright © 2016 pulupulu. All rights reserved.
//

import UIKit

class StatTableViewCell: UITableViewCell {

	private var nameLabel: StatLabel?
	private var levelLabel: StatLabel?
	private var gradeLabel: StatLabel?
	private var goldLabel: StatLabel?
	private var energyIcon: UIImageView?
	private var energyLabel: StatLabel?
	private var comboLabel: StatLabel?

	required init?(coder aDecoder: NSCoder) {
		super.init(coder: aDecoder)
		commonInit()
	}

	override init(style: UITableViewCellStyle, reuseIdentifier: String?) {
		super.init(style: style, reuseIdentifier: reuseIdentifier)
		commonInit()
	}

	func commonInit() {
		nameLabel = StatLabel()
		levelLabel = StatLabel()
		gradeLabel = StatLabel()
	}

	func renderData(data: [String: Any]) {
	}
}

private class StatLabel: UILabel {
	required init?(coder aDecoder: NSCoder) {
		super.init(coder: aDecoder)
		commonInit()
	}
	override init(frame: CGRect) {
		super.init(frame: frame)
		commonInit()
	}

	func commonInit() {
		self.font = UIFont(name: PLConstants.usualFont, size: 25)
		self.textColor = UIColor.whiteColor()
	}
}
