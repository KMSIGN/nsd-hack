import 'dart:io';
import 'dart:typed_data';

import 'package:crypto/crypto.dart';
import 'package:file_handler/utils.dart';

const _blockSize = 8 * 1024 * 1024; //8Mb

class FileHashMap {
  final File file;
  List<Digest> hashes = [];

  FileHashMap(this.file);

  Future<List<Digest>> parseFileToHashMap(Progress progress) async {
    //final chunkedFile = asChunkedStream<int>(_blockSize, file.openRead().expand((e) => e));
    //return chunkedFile.map(sha1.convert).toList();

    final rafile = await file.open();
    final buffer = Uint8List(_blockSize);
    int bufferLen;
    progress.total = (file.lengthSync() / _blockSize).ceil();

    do {
      bufferLen = await rafile.readInto(buffer, 0);
      hashes.add(sha1.convert(buffer));
      progress.add();
    } while (bufferLen == _blockSize);
    return hashes;
  }
}
