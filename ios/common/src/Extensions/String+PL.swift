//
//  String+PL.swift
//  admin
//
//  Created by tassar on 5/2/16.
//  Copyright © 2016 pulupulu. All rights reserved.
//

import Foundation
extension String {

	subscript(i: Int) -> Character {
		return self[self.characters.index(self.startIndex, offsetBy: i)]
	}

	subscript(i: Int) -> String {
		return String(self[i] as Character)
	}

	subscript(r: Range<Int>) -> String {
		let start = characters.index(startIndex, offsetBy: r.lowerBound)
        let end = characters.index(start, offsetBy: r.upperBound - r.lowerBound)
		return self[Range(start ..< end)]
	}

	func last() -> String {
		return String(self[self.characters.index(before: self.endIndex)])
	}
}
