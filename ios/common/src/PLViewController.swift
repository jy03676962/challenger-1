//
//  PLViewController.swift
//  admin
//
//  Created by tassar on 4/23/16.
//  Copyright © 2016 pulupulu. All rights reserved.
//

import UIKit
import EasyPeasy

class PLViewController: UIViewController {
	override func viewDidLoad() {
		super.viewDidLoad()
		let imageView = UIImageView()
		imageView.image = UIImage(named: "GlobalBackground")
		view.insertSubview(imageView, atIndex: 0)
		imageView <- Edges()
	}
	override func prefersStatusBarHidden() -> Bool {
		return true
	}
}
