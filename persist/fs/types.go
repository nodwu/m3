	"github.com/m3db/m3db/storage/namespace"
	xtime "github.com/m3db/m3x/time"
	Open(namespace ts.ID, blockSize time.Duration, shard uint32, start time.Time) error
	// ReadBloomFilter returns the bloom filter stored on disk in a container object that is safe
	// for concurrent use and has a Close() method for releasing resources when done.
	ReadBloomFilter() (*ManagedConcurrentBloomFilter, error)

	Open(md namespace.Metadata) error
	Open(md namespace.Metadata) error
	// Validate will validate the options and return an error if not valid
	Validate() error

	// SetIndexSummariesPercent size sets the percent of index summaries to write
	SetIndexSummariesPercent(value float64) Options

	// IndexSummariesPercent size returns the percent of index summaries to write
	IndexSummariesPercent() float64

	// SetIndexBloomFilterFalsePositivePercent size sets the percent of false positive
	// rate to use for the index bloom filter size and k hashes estimation
	SetIndexBloomFilterFalsePositivePercent(value float64) Options

	// IndexBloomFilterFalsePositivePercent size returns the percent of false positive
	// rate to use for the index bloom filter size and k hashes estimation
	IndexBloomFilterFalsePositivePercent() float64

	// SetInfoReaderBufferSize sets the buffer size for reading TSDB info, digest and checkpoint files
	SetInfoReaderBufferSize(value int) Options

	// InfoReaderBufferSize returns the buffer size for reading TSDB info, digest and checkpoint files
	InfoReaderBufferSize() int

	// SetDataReaderBufferSize sets the buffer size for reading TSDB data and index files
	SetDataReaderBufferSize(value int) Options

	// DataReaderBufferSize returns the buffer size for reading TSDB data and index files
	DataReaderBufferSize() int

	// SetSeekReaderBufferSize size sets the buffer size for seeking TSDB files
	SetSeekReaderBufferSize(value int) Options

	// SeekReaderBufferSize size returns the buffer size for seeking TSDB files
	SeekReaderBufferSize() int

	// SetMmapEnableHugePages sets whether mmap huge pages are enabled when running on linux
	SetMmapEnableHugePages(value bool) Options

	// MmapEnableHugePages returns whether mmap huge pages are enabled when running on linux
	MmapEnableHugePages() bool

	// SetMmapHugePagesThreshold sets the threshold when to use mmap huge pages for mmap'd files on linux
	SetMmapHugePagesThreshold(value int64) Options
	// MmapHugePagesThreshold returns the threshold when to use mmap huge pages for mmap'd files on linux
	MmapHugePagesThreshold() int64