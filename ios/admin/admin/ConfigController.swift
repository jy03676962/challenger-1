//
//  ConfigController.swift
//  admin
//
//  Created by tassar on 4/23/16.
//  Copyright © 2016 pulupulu. All rights reserved.
//

import UIKit
import AutoKeyboardScrollView
import EasyPeasy
import SwiftyUserDefaults

class ConfigController: PLViewController {

	@IBOutlet weak var wrapperView: UIView!
	@IBOutlet weak var hostTextField: UITextField!

	@IBAction func saveConfig() {
		if let text = hostTextField.text {
			Defaults[.host] = text
			WsClient.singleton.connect(PLConstants.getWsAddress())
		}
	}

	override func viewDidLoad() {
		super.viewDidLoad()
		let scrollView = AutoKeyboardScrollView()
		scrollView.backgroundColor = UIColor.clearColor()
		view.addSubview(scrollView)
		wrapperView.removeFromSuperview()
		scrollView.addSubview(wrapperView)
		scrollView <- Edges()
		wrapperView <- Edges()
	}
}
