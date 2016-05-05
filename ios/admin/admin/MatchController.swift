//
//  SecondViewController.swift
//  admin
//
//  Created by tassar on 4/20/16.
//  Copyright © 2016 pulupulu. All rights reserved.
//

import UIKit
import Alamofire
import AlamofireImage
import EasyPeasy
import ObjectMapper

let cellSize = 45
let cellBorder = 10

class MatchController: PLViewController {
	@IBOutlet weak var groupIDLabel: UILabel!
	@IBOutlet weak var matchStatusLabel: UILabel!
	@IBOutlet weak var playerCountLabel: UILabel!
	@IBOutlet weak var totalCoinLabel: UILabel!
	@IBOutlet weak var energyLabel: UILabel!
	@IBOutlet weak var matchTimeLabel: UILabel!
	@IBOutlet weak var matchModeImageView: UIImageView!
	@IBOutlet weak var mapContainerView: UIView!
	@IBOutlet weak var playerTableView: UITableView!

	var match: Match?

	var mapView: UIImageView = UIImageView()

	@IBAction func forceEnd() {
	}

	override func viewDidLoad() {
		super.viewDidLoad()
	}

	override func viewWillAppear(animated: Bool) {
		super.viewWillAppear(animated)
		DataManager.singleton.subscribeData([.UpdateMatch], receiver: self)
		if mapView.image == nil {
			Alamofire.request(.GET, PLConstants.getHttpAddress("api/asset/map.png"))
				.responseImage(completionHandler: { response in
					if let image = response.result.value {
						self.mapView.image = image
						self.mapContainerView.addSubview(self.mapView)
						self.mapView <- [
							Size(image.size),
							Center()
						]
					}
			})
		}
	}

	override func viewDidDisappear(animated: Bool) {
		super.viewDidDisappear(animated)
		DataManager.singleton.unsubscribe(self)
	}

	func renderMatch() {
	}
}

extension MatchController: DataReceiver {
	func onReceivedData(json: [String: AnyObject], type: DataType) {
		if type == .UpdateMatch {
			match = Mapper<Match>().map(json["data"] as! String)
		}
	}
}

extension MatchController: UITableViewDelegate, UITableViewDataSource {
	func tableView(tableView: UITableView, numberOfRowsInSection section: Int) -> Int {
		return 0
	}

	func numberOfSectionsInTableView(tableView: UITableView) -> Int {
		return 1
	}

	func tableView(tableView: UITableView, cellForRowAtIndexPath indexPath: NSIndexPath) -> UITableViewCell {
		return UITableViewCell()
	}
}
