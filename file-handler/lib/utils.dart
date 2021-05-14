typedef ProgressCallback = void Function(Progress progress);

class Progress {
  int total;
  int current;
  ProgressCallback? callback;

  DateTime timeStart = DateTime.now();

  Duration get timeit => timeStart.difference(DateTime.now());

  double get percent => (current / total) * 100;

  Progress({this.callback, this.total = 100, this.current = 0});

  void add([int? count]) {
    current += count ?? 1;
    callback?.call(this);
  }
}
