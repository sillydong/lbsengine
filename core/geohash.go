package core

import "github.com/mmcloughlin/geohash"

func LoopNeighbours(latitude, longitude float64, precision uint, loop int) (neighbours []string) {
	if loop == 0 {
		loop = 1
	}
	hash := geohash.EncodeWithPrecision(latitude, longitude, precision)
	if loop == 1 {
		neighbours = geohash.Neighbors(hash)
		neighbours = append(neighbours, hash)
		return neighbours
	} else {
		neighbours = append(neighbours, hash)

		box := geohash.BoundingBox(hash)
		centerlat, centerlng := box.Center()
		height := box.MaxLat - box.MinLat
		width := box.MaxLng - box.MinLng

		for i := 1; i <= loop; i++ {
			side := 2 * i
			fi := float64(i)
			latup := centerlat + width*fi
			lngright := centerlng + width*fi
			latdown := centerlat - height*fi
			lngleft := centerlng - width*fi

			for k := 0; k < side; k++ {
				//上
				uplng := centerlng + width*float64(k-i+1)
				neighbours = append(neighbours, geohash.EncodeWithPrecision(latup, uplng, precision))

				//右
				rightlat := centerlat - height*float64(k-i+1)
				neighbours = append(neighbours, geohash.EncodeWithPrecision(rightlat, lngright, precision))

				//下
				downlng := centerlng - width*float64(k-i+1)
				neighbours = append(neighbours, geohash.EncodeWithPrecision(latdown, downlng, precision))

				//左
				leftlat := centerlat + height*float64(k-i+1)
				neighbours = append(neighbours, geohash.EncodeWithPrecision(leftlat, lngleft, precision))
			}
		}

		return neighbours
	}
}
