//
//  HistoryController.swift
//  admin
//
//  Created by tassar on 4/20/16.
//  Copyright © 2016 pulupulu. All rights reserved.
//

import UIKit
import Alamofire
import AlamofireObjectMapper
import PKHUD

class HistoryController: PLViewController {
	@IBOutlet weak var tableView: UITableView!
	@IBOutlet weak var startAnswerButton: UIButton!

	var refreshControl: UIRefreshControl!
	var data: [MatchData]?

	@IBAction func startAnswer() {
		HUD.show(.Progress)
		Alamofire.request(.POST, PLConstants.getHttpAddress("api/start_answer"))
			.responseObject(completionHandler: { (response: Response<MatchData, NSError>) in
				let md = response.result.value
				// TODO
				if self.data != nil && md != nil {
					for (i, d) in self.data!.enumerate() {
						if d.id == md!.id {
							self.data![i] = md!
							break
						}
					}
				}
				HUD.hide()
		})
		performSegueWithIdentifier("PresentMatchResult", sender: nil)
	}
	override func viewDidLoad() {
		super.viewDidLoad()
		refreshControl = UIRefreshControl()
		refreshControl.addTarget(self, action: #selector(refreshHistory), forControlEvents: .ValueChanged)
		tableView.addSubview(refreshControl)
	}
	override func viewWillAppear(animated: Bool) {
		super.viewWillAppear(animated)
		if data == nil {
			refreshHistory()
		}
	}

	func refreshHistory() {
		Alamofire.request(.GET, PLConstants.getHttpAddress("api/history"))
			.responseArray(completionHandler: { (response: Response<[MatchData], NSError>) in
				self.data = response.result.value
				if self.data != nil {
					self.tableView.reloadData()
				}
				self.refreshControl.endRefreshing()
		})
	}
}

extension HistoryController: UITableViewDelegate, UITableViewDataSource {
	func tableView(tableView: UITableView, numberOfRowsInSection section: Int) -> Int {
		return self.data == nil ? 0 : self.data!.count
	}
	func numberOfSectionsInTableView(tableView: UITableView) -> Int {
		return 1
	}
	func tableView(tableView: UITableView, cellForRowAtIndexPath indexPath: NSIndexPath) -> UITableViewCell {
		let cell = tableView.dequeueReusableCellWithIdentifier("HistoryTableViewCell") as! HistoryTableViewCell
		cell.setData(data![indexPath.row])
		return cell
	}
	func tableView(tableView: UITableView, didSelectRowAtIndexPath indexPath: NSIndexPath) {
		startAnswerButton.enabled = true
	}
	func tableView(tableView: UITableView, didDeselectRowAtIndexPath indexPath: NSIndexPath) {
		startAnswerButton.enabled = false
	}
}
