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

	override func viewDidLoad() {
		super.viewDidLoad()
		if let image = self.backgroundImage() {
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

	func backgroundImage() -> UIImage? {
		return UIImage(named: "GlobalBackground")
	}
}
