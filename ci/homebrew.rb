class DockerCommander < Formula
  desc "Execute commands in docker containers."
  homepage "https://github.com/daylioti/docker-commander"
  url "https://github.com/daylioti/docker-commander.git",
    :tag      => "1.1.4",
    :revision => "76a727fc77a4236a1c980f6b637ffa7288816c8b"

  version "1.1.4"

  depends_on "go" => :build

  def install
    ENV["GOPATH"] = buildpath
    ENV["GO111MODULE"] = "on"
    src = buildpath/"src/github.com/daylioti/docker-commander"
    src.install buildpath.children
    src.cd do
      system "make", "build"
      bin.install "docker-commander"
      prefix.install_metafiles
    end
  end

  test do
      system "#{bin}/docker-commander", "-v"
  end

end