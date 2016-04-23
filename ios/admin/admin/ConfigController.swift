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

class ConfigController: UIViewController {

	@IBOutlet weak var wrapperView: UIView!
	@IBOutlet weak var hostTextField: UITextField!

	@IBAction func saveConfig() {
		if let text = hostTextField.text {
			Defaults[.host] = text
			NSNotificationCenter.defaultCenter().postNotificationKey(.HostChanged, object: nil)
		}
	}

	override func viewDidLoad() {
		super.viewDidLoad()
		let scrollView = AutoKeyboardScrollView()
		view.addSubview(scrollView)
		wrapperView.removeFromSuperview()
		scrollView.addSubview(wrapperView)
		scrollView <- Edges()
		wrapperView <- Edges()
	}

	override func prefersStatusBarHidden() -> Bool {
		return true
	}
}
