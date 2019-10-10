#!/bin/bash

# This is a slightly extended version of https://github.com/lquesada/voxel-3d-models/tree/master/tools/voxobjrename
# Takes as input .vox, .mtl, .obj, and .png files from the sources directory, renames them into the assets directory
# then converts them into Go source code so they're embedded in the binary, and cleans up "assets".
#
# Read the documentation at the github page for requirements, more details, and a less hacky version of this script.

(cd assets
find -name \*.obj -exec rm '{}' \;
find -name \*.mtl -exec rm '{}' \;
find -name \*.png -exec rm '{}' \;
)
rm -rf extract
cp -rf sources extract
(cd extract

af='../assets/resources.go'
rm -f $af
echo 'package assets' > $af
echo '' >> $af
echo 'var Files map[string][]byte = map[string][]byte{' >> $af

for x in `ls *.vox`;do
  file=`echo $x | sed s/'\.vox$'//g`
  i=0
  for f in ` \
strings $file.vox | \
grep -e VGCX | \
sed s/'^.*GCX_\(.*\)_VGC.*$'/'\1'/g | \
sed s/'^'/'assets\/'/g \
`;do
    echo $f | grep _VGIGNOREVG_
    retVal=$?
    if [ $retVal -eq 0 ]; then
      rm $file-$i.{obj,mtl,png}
      i=$((i+1))
      continue
    fi

    for y in obj mtl png;do
      echo "$file-$i.$y --> $f.$y"
      mv $file-$i.$y ../$f.$y
    done
    b=$(basename $f)
    sed -i -- s/"$file-$i"/"$b"/g ../$f.obj
    sed -i -- 's/\r$//' ../$f.obj
    sed -i -- s/"$file-$i"/"$b"/g ../$f.mtl
    sed -i -- 's/\r$//' ../$f.mtl
    sed -i -- 's/illum 1/illum 1/g' ../$f.mtl
    sed -i -- 's/Ka 0.000 0.000 0.000/Ka 1.000 1.000 1.000/g' ../$f.mtl
    sed -i -- 's/Kd 1.000 1.000 1.000/Kd 1.000 1.000 1.000/g' ../$f.mtl
    sed -i -- 's/Ks 0.000 0.000 0.000/Ks 0.000 0.000 0.000/g' ../$f.mtl
    grep -q ../$f.mtl -e "^Ns 200$" || echo "Ns 200" >> ../$f.mtl
    grep -q ../$f.mtl -e "^d 1$" || echo "d 1" >> ../$f.mtl
    for y in obj mtl png;do
      (echo -n "\"$f.$y\": []byte{"; xxd -i ../$f.$y | sed -e '$ d' -e '1d' | sed -e '$ d' | tr -d '\n' | sed s/'  '/' '/g | sed s/'^ '/''/g;echo '},') >> $af
    done
    i=$((i+1))
  done
  rm $x
done

echo '}' >> $af

cd ..
echo "Files left?"
ls extract && rm -rf extract
)

(cd assets
find -name \*.obj -exec rm '{}' \;
find -name \*.mtl -exec rm '{}' \;
find -name \*.png -exec rm '{}' \;
)
