import sys

from pyproj import Proj, transform

with open(sys.argv[1], 'w') as output:
    print >>output, 'package proj'
    print >>output
    print >>output, '//go:generate python generate-testdata.py testdata_test.go'
    print >>output
    print >>output, 'var ('
    print >>output, '\tepsg3857TestData = []struct {'
    print >>output, '\t\tlon, lat float64'
    print >>output, '\t\te, n     float64'
    print >>output, '\t}{'
    p1 = Proj(proj='latlong')
    p2 = Proj(init='epsg:3857')
    for lon in xrange(-180, 180, 30):
        for lat in xrange(-80, 90, 30):
            x, y = transform(p1, p2, lon, lat)
            print >>output, '\t\t{lat: rad(%d), lon: rad(%d), e: %f, n: %f},' % (lat, lon, x, y)
    print >>output, '\t}'
    print >>output, ')'
