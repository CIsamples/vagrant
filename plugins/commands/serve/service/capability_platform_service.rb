require "google/protobuf/well_known_types"

module VagrantPlugins
  module CommandServe
    module Service
      module CapabilityPlatformService

        def self.included(klass)
          klass.include(Util::ServiceInfo)
          klass.prepend(Util::HasMapper)
          klass.prepend(Util::HasBroker)
          klass.prepend(Util::HasLogger)
          klass.prepend(Util::ExceptionLogger)

          klass.class_eval do
            attr_reader :capabilities, :default_args
          end
        end

        def initialize_capability_platform!(capabilities, default_args)
          @capabilities = capabilities
          @default_args = default_args
        end

        # TODO(spox): request scoping for seed values - needs investigation
        def seed(req, ctx)
          @seeds = req.arguments.to_a
          Empty.new
        end

        def seeds(req, ctx)
          SDK::Args::Direct.new(arguments: @seeds)
        end

        def has_capability_spec(*_)
          SDK::FuncSpec.new(
            name: "has_capability_spec",
            args: [
              SDK::FuncSpec::Value.new(
                type: "hashicorp.vagrant.sdk.Args.NamedCapability",
                name: "",
              )
            ],
            result: [
              SDK::FuncSpec::Value.new(
                type: "hashicorp.vagrant.sdk.Platform.Capability.CheckResp",
                name: "",
              ),
            ],
          )
        end

        def has_capability(req, ctx)
          with_info(ctx) do |info|
            cap_name = mapper.funcspec_map(req)
            plugin_name = info.plugin_name
            logger.debug("checking for #{cap_name} capability in #{plugin_name}")

            caps_registry = @capabilities[plugin_name.to_sym]
            has_cap = caps_registry.key?(cap_name.to_sym)

            SDK::Platform::Capability::CheckResp.new(
              has_capability: has_cap
            )
          end
        end

        def capability_spec(req, ctx)
          SDK::FuncSpec.new(
            name: "capability_spec",
            args: default_args.values + [
              SDK::FuncSpec::Value.new(
                type: "hashicorp.vagrant.sdk.Args.Direct",
                name: "",
              )
            ],
            result: [
              SDK::FuncSpec::Value.new(
                type: "hashicorp.vagrant.sdk.Platform.Capability.Resp",
                name: "",
              )
            ]
          )
        end

        def capability(req, ctx)
          with_info(ctx) do |info|
            cap_name = req.name.to_sym
            plugin_name = info.plugin_name.to_sym

            logger.debug("executing capability #{cap_name} on plugin #{plugin_name}")

            caps_registry = capabilities[plugin_name]
            target_cap = caps_registry.get(cap_name)

            args = mapper.funcspec_map(
              req.func_args,
              expect: default_args.keys + [Types::Direct]
            )
            args = capability_arguments(args)
            cap_method = target_cap.method(cap_name)

            arg_list = args.join("\n  - ")
            logger.debug("arguments to be passed to #{cap_name} on plugin #{plugin_name}:\n  - #{arg_list}")

            result = cap_method.call(*args)

            val = Google::Protobuf::Value.new
            val.from_ruby(result)
            SDK::Platform::Capability::Resp.new(
              result: Google::Protobuf::Any.pack(val)
            )
          end
        end

        def capability_arguments(args)
          direct = args.pop
          args + direct.arguments
        end
      end
    end
  end
end
