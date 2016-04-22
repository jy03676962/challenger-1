//
//  FirstViewController.swift
//  admin
//
//  Created by tassar on 4/20/16.
//  Copyright © 2016 pulupulu. All rights reserved.
//

import UIKit

class HallController: UIViewController {

	override func viewDidLoad() {
		super.viewDidLoad()
		setBackgroundImage("GlobalBackground")
	}
	override func prefersStatusBarHidden() -> Bool {
		return true
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
		return UITableViewCell()
	}
}
