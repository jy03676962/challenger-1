//
//  HallTableViewCell.swift
//  admin
//
//  Created by tassar on 4/23/16.
//  Copyright © 2016 pulupulu. All rights reserved.
//

import UIKit
import SwiftyJSON
import SWTableViewCell

class HallTableViewCell: SWTableViewCell {

	@IBOutlet weak var teamIDLabel: UILabel!
	@IBOutlet weak var teamSizeLabel: UILabel!
	@IBOutlet weak var delayCountLabel: UILabel!
	@IBOutlet weak var waitTimeLabel: UILabel!
	@IBOutlet weak var delayCountImageView: UIImageView!
	@IBOutlet weak var numberLabel: UILabel!
	override func awakeFromNib() {
		super.awakeFromNib()
		backgroundColor = UIColor.clearColor()
	}

	override func setSelected(selected: Bool, animated: Bool) {
		super.setSelected(selected, animated: animated)

		// Configure the view for the selected state
	}

	func setData(dict: JSON, number: Int) {
		teamIDLabel.text = dict["id"].stringValue
		teamSizeLabel.text = String(dict["size"].intValue)
		let delayCount = dict["delayCount"].intValue
		delayCountLabel.text = "- \(delayCount) -"
		var waitTime = dict["waitTime"].intValue
		let waitHour = waitTime / 3600
		waitTime -= 3600 * waitHour
		let waitMin = waitTime / 60
		waitTime -= 60 * waitMin
		waitTimeLabel.text = String(format: "%02d:%02d:%02d", waitHour, waitMin, waitTime)
		let delayImageName = "IconLate\(delayCount)"
		let delayImage = UIImage(named: delayImageName)
		delayCountImageView.image = delayImage
		numberLabel.text = String(number + 1)
	}
}