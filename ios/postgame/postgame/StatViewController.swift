//
//  ViewController.swift
//  postgame
//
//  Created by tassar on 3/29/16.
//  Copyright © 2016 pulupulu. All rights reserved.
//

import UIKit

class StatViewController: PLBaseViewController {

	override func viewDidLoad() {
		super.viewDidLoad()
		// Do any additional setup after loading the view, typically from a nib.
	}

	override func viewWillAppear(animated: Bool) {
		super.viewWillAppear(animated)
		navigationController?.navigationBar.hidden = true
	}

	override func backgroundImage() -> UIImage? {
		return UIImage(named: "GlobalBackground")
	}

	override func prefersStatusBarHidden() -> Bool {
		return true
	}
}
