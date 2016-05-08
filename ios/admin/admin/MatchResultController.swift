//
//  MatchResultController.swift
//  admin
//
//  Created by tassar on 5/6/16.
//  Copyright © 2016 pulupulu. All rights reserved.
//

import UIKit

class MatchResultController: PLViewController {
	var matchData: MatchData?

	@IBOutlet weak var headerImageView: UIImageView!
	@IBOutlet weak var teamIDLabel: UILabel!
	@IBOutlet weak var scoreLabel: UILabel!

	@IBAction func endAnswer() {
	}

	override func viewDidLoad() {
		super.viewDidLoad()
		renderData()
	}

	func renderData() {
		if let data = matchData {
			let image = data.mode == "g" ? UIImage(named: "FunImage") : UIImage(named: "SurivalImage")
			headerImageView.image = image
			teamIDLabel.text = data.teamID
		}
	}
}

extension MatchResultController: UITableViewDataSource, UITableViewDelegate {
	func numberOfSectionsInTableView(tableView: UITableView) -> Int {
		return 1
	}

	func tableView(tableView: UITableView, numberOfRowsInSection section: Int) -> Int {
		return matchData == nil ? 0 : matchData!.member.count
	}

	func tableView(tableView: UITableView, cellForRowAtIndexPath indexPath: NSIndexPath) -> UITableViewCell {
		let cell = tableView.dequeueReusableCellWithIdentifier("MatchResultCell") as! MatchResultCell
		let data = matchData!.member[indexPath.row]
		cell.setData(data)
		return cell
	}
}
