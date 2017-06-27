from pyproj import Proj

print 'package proj'
print
print 'var ('
print '\tepsg3857TestData = []struct{'
print '\t\tlon, lat float64'
print '\t\te, n     float64'
print '\t}{'
epsg3857 = Proj(init='epsg:3857')
for lon in xrange(-170, 180, 10):
    for lat in xrange(-80, 90, 10):
        e, n = epsg3857(lat, lon)
        print '\t\t{lon: %d, lat: %d, e: %f, n: %f},' % (lon, lat, e, n)
print '\t}'
print ')'
