//
//  PLBaseViewController.swift
//  postgame
//
//  Created by tassar on 3/31/16.
//  Copyright © 2016 pulupulu. All rights reserved.
//

import UIKit
import Masonry
import HEXColor

class PLBaseViewController: UIViewController {

	var backgroundImageName: String = "GlobalBackground"

	override func viewDidLoad() {
		super.viewDidLoad()
		if let image = UIImage(named: backgroundImageName) {
			let imageView = UIImageView()
			imageView.image = image
			view.insertSubview(imageView, atIndex: 0)
			imageView.mas_makeConstraints { make in
				make.edges.equalTo()(self.view)
			}
		}
	}

	override func prefersStatusBarHidden() -> Bool {
		return true
	}
}
