//
//  HallController.swift
//  admin
//
//  Created by tassar on 4/20/16.
//  Copyright © 2016 pulupulu. All rights reserved.
//

import UIKit
import UIScrollView_InfiniteScroll

class HallController: PLViewController {
	@IBOutlet weak var teamtableView: UITableView!
	@IBOutlet weak var teamIDLabel: UILabel!

	@IBOutlet weak var modeImageView: UIImageView!
	@IBOutlet weak var modeLabel: UILabel!
	@IBOutlet weak var playerNumberLabel: UILabel!
	@IBOutlet weak var readyButton: UIButton!
	@IBOutlet weak var startButton: UIButton!
	var refreshControl: UIRefreshControl!

	override func viewDidLoad() {
		super.viewDidLoad()
		refreshControl = UIRefreshControl()
		refreshControl.addTarget(self, action: #selector(HallController.refreshTeamData), forControlEvents: UIControlEvents.ValueChanged)
		teamtableView.addSubview(refreshControl)
		teamtableView.addInfiniteScrollWithHandler({ scrollView in
			let tableView = scrollView as! UITableView
			tableView.finishInfiniteScroll()
		})
	}
	override func viewWillAppear(animated: Bool) {
	}
	func refreshTeamData() {
		refreshControl.endRefreshing()
	}
	@IBAction func changeMode(sender: UITapGestureRecognizer) {
	}
	@IBAction func callTeam(sender: UIButton) {
	}
	@IBAction func delayTeam(sender: UIButton) {
	}
	@IBAction func addPlayer(sender: UIButton) {
	}
	@IBAction func removePlayer(sender: UIButton) {
	}
	@IBAction func ready(sender: UIButton) {
	}
	@IBAction func start(sender: AnyObject) {
	}
}

extension HallController: UITableViewDataSource, UITableViewDelegate {
	func tableView(tableView: UITableView, numberOfRowsInSection section: Int) -> Int {
		return 1
	}
	func numberOfSectionsInTableView(tableView: UITableView) -> Int {
		return 1
	}
	func tableView(tableView: UITableView, cellForRowAtIndexPath indexPath: NSIndexPath) -> UITableViewCell {
		return tableView.dequeueReusableCellWithIdentifier("HallTableViewCell")!
	}
}
