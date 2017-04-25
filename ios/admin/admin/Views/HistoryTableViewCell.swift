//
//  HistoryTableViewCell.swift
//  admin
//
//  Created by tassar on 5/6/16.
//  Copyright © 2016 pulupulu. All rights reserved.
//

import UIKit

class HistoryTableViewCell: UITableViewCell {
	@IBOutlet weak var teamIDLabel: UILabel!
	@IBOutlet weak var playerCountLabel: UILabel!
	@IBOutlet weak var statusLabel: UILabel!
	@IBOutlet weak var backgroundImageView: UIImageView!
	override func awakeFromNib() {
		super.awakeFromNib()
		backgroundColor = UIColor.clear
	}
	func setData(_ data: MatchData) {
		teamIDLabel.text = data.teamID
		playerCountLabel.text = "\(data.member.count)"
		var txt: String
		switch data.answerType! {
		case .notAnswer:
			txt = "尚未答题"
		case .answering:
			txt = "答题中"
		case .answered:
			txt = "已答题"
		}
		statusLabel.text = txt
	}
	override func setSelected(_ selected: Bool, animated: Bool) {
		backgroundImageView.isHidden = !selected
	}
}
