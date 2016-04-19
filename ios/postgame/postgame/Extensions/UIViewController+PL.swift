//
//  UIViewController+PL.swift
//  postgame
//
//  Created by tassar on 4/19/16.
//  Copyright © 2016 pulupulu. All rights reserved.
//

import Foundation
import UIKit
/*
 Pulupulu iOS项目中UIViewController的常用方法
 */
extension UIViewController {
	func setBackgroundImage(imageName: String) {
		if let image = UIImage(named: imageName) {
			var imageView: UIImageView? = view.viewWithTag(ReservedViewTag.vc_BackgroundImageView) as? UIImageView
			if imageView == nil {
				imageView = UIImageView()
				imageView!.tag = ReservedViewTag.vc_BackgroundImageView
			}
			imageView!.image = image
			view.insertSubview(imageView!, atIndex: 0)
			imageView!.mas_makeConstraints { make in
				make.edges.equalTo()(self.view)
			}
		}
	}
}