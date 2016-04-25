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
	override func awakeFromNib() {
		super.awakeFromNib()
		backgroundColor = UIColor.clearColor()
	}

	override func setSelected(selected: Bool, animated: Bool) {
		super.setSelected(selected, animated: animated)

		// Configure the view for the selected state
	}

	func setData(dict: JSON) {
		teamIDLabel.text = dict["id"].stringValue
	}
}