# Homebrew formula for garfield-router
# This is a template; goreleaser generates the actual formula in the tap repo.
# Manual install: brew install garfield/tap/garfield-router

class GarfieldRouter < Formula
  desc "Intelligent LLM proxy with cost-optimised routing"
  homepage "https://github.com/Garfield/garfield-router"
  license "Apache-2.0"
  version "1.0.0"

  on_macos do
    on_arm do
      url "https://github.com/Garfield/garfield-router/releases/download/v#{version}/garfield-router_#{version}_darwin_arm64.tar.gz"
      # sha256 will be filled by goreleaser
    end
    on_intel do
      url "https://github.com/Garfield/garfield-router/releases/download/v#{version}/garfield-router_#{version}_darwin_amd64.tar.gz"
    end
  end

  on_linux do
    on_arm do
      url "https://github.com/Garfield/garfield-router/releases/download/v#{version}/garfield-router_#{version}_linux_arm64.tar.gz"
    end
    on_intel do
      url "https://github.com/Garfield/garfield-router/releases/download/v#{version}/garfield-router_#{version}_linux_amd64.tar.gz"
    end
  end

  def install
    bin.install "garfield-router"
    etc.install "config.yaml" => "garfield-router/config.yaml"
  end

  def post_install
    ohai "Run 'garfield-router init' to auto-detect backends and generate config"
    ohai "Run 'garfield-router' to start the router"
    ohai "Dashboard at http://localhost:8050"
  end

  service do
    run [opt_bin/"garfield-router"]
    keep_alive true
    working_dir var/"garfield-router"
    log_path var/"log/garfield-router.log"
    error_log_path var/"log/garfield-router.log"
    environment_variables CONFIG_PATH: etc/"garfield-router/config.yaml"
  end

  test do
    assert_match "garfield-router v", shell_output("#{bin}/garfield-router version")
  end
end
