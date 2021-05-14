import 'dart:io';

import 'package:file_handler/file_hashmap.dart';
import 'package:file_handler/utils.dart';

void main(List<String> arguments) async {
  final progress = Progress(
    callback: (progress) => print(progress.percent),
  );

  final hashmap = await FileHashMap(File('/home/royalcat/neiro.7z')).parseFileToHashMap(progress);
  print(progress.timeit);
}
